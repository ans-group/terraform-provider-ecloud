package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceSolutionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolutionTemplateCreate,
		Read:   resourceSolutionTemplateRead,
		Update: resourceSolutionTemplateUpdate,
		Delete: resourceSolutionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"solution_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"virtualmachine_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSolutionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	solutionID := d.Get("solution_id").(int)
	vmID := d.Get("virtualmachine_id").(int)
	templateName := d.Get("name").(string)

	createReq := ecloudservice.CreateVirtualMachineTemplateRequest{
		TemplateName: templateName,
		TemplateType: ecloudservice.TemplateTypeSolution,
	}
	log.Printf("Created CreateVirtualMachineTemplateRequest: %+v", createReq)

	log.Printf("Retrieving virtual machine with id [%d]", vmID)
	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("Error retrieving virtual machine: %s", err)
	}

	if solutionID != vm.SolutionID {
		return fmt.Errorf("Invalid solution id [%d], expected [%d]", solutionID, vm.SolutionID)
	}

	log.Print("Creating solution template")
	err = service.CreateVirtualMachineTemplate(vmID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating solution template: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Existent"},
		Refresh:    SolutionTemplateExistentStatusRefreshFunc(service, solutionID, templateName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for solution template with name [%s] to return status of [Existent]: %s", templateName, err)
	}

	d.SetId(templateName)

	return resourceSolutionTemplateRead(d, meta)
}

func resourceSolutionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	templateName := d.Id()
	solutionID := d.Get("solution_id").(int)

	log.Printf("Retrieving template with name [%s] for solution with ID [%d]", templateName, solutionID)
	template, err := service.GetSolutionTemplate(solutionID, templateName)
	if err != nil {
		switch err.(type) {
		case *ecloudservice.TemplateNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", template.Name)
	d.Set("solution_id", solutionID)

	return nil
}

func resourceSolutionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	templateName := d.Id()
	solutionID := d.Get("solution_id").(int)

	if d.HasChange("name") {
		newTemplateName := d.Get("name").(string)

		err := service.RenameSolutionTemplate(solutionID, templateName, ecloudservice.RenameTemplateRequest{
			Destination: newTemplateName,
		})
		if err != nil {
			return fmt.Errorf("Error updating template with name [%s]: %s", templateName, err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{"Existent"},
			Refresh:    SolutionTemplateExistentStatusRefreshFunc(service, solutionID, newTemplateName),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for solution template with name [%s] to return status of [Existent]: %s", newTemplateName, err)
		}

		d.SetId(newTemplateName)
	}

	return resourceSolutionTemplateRead(d, meta)
}

func resourceSolutionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	templateName := d.Id()
	solutionID := d.Get("solution_id").(int)

	log.Printf("Removing solution template with name [%s] for solution with ID [%d]", templateName, solutionID)
	err := service.DeleteSolutionTemplate(solutionID, templateName)
	if err != nil {
		return fmt.Errorf("Error removing solution template with name [%s] for solution with ID [%d]: %s", templateName, solutionID, err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"NonExistent"},
		Refresh:    SolutionTemplateExistentStatusRefreshFunc(service, solutionID, templateName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for solution template with name [%s] to return status of [NonExistent]: %s", templateName, err)
	}

	return nil
}

// SolutionTemplateExistentStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func SolutionTemplateExistentStatusRefreshFunc(service ecloudservice.ECloudService, solutionID int, templateName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		template, err := service.GetSolutionTemplate(solutionID, templateName)
		if err != nil {
			if _, ok := err.(*ecloudservice.TemplateNotFoundError); ok {
				return template, "NonExistent", nil
			}
			return nil, "", err
		}

		return template, "Existent", nil
	}
}

// func getSolutionTemplateResourceID(solutionID int, templateName string) string {
// 	return fmt.Sprintf("%d::%s", solutionID, templateName)
// }

// func parseSolutionTemplateResourceID(id string) (solutionID int, templateName string) {
// 	parts := strings.Split(id, "::")
// 	solutionID, _ = strconv.Atoi(parts[0])
// 	return solutionID, parts[1]
// }
