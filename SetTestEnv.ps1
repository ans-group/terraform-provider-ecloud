Param
(
    $Force=$false
)

$env:TF_ACC = "1"
function SetEnv($Name)
{
    if (((Test-Path env:\$Name) -ne $true) -or $Force)
    {
        New-Item -Name $Name -value (Read-Host -Prompt "Enter $Name") -ItemType Variable -Path Env: | Out-Null
    }
}

SetEnv "UKF_API_KEY"
SetEnv "UKF_TEST_VPC_ID"
SetEnv "UKF_TEST_DHCP_AVAILABILITY_ZONE_ID"