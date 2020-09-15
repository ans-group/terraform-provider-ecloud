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
SetEnv "UKF_TEST_VPC_NAME"
SetEnv "UKF_TEST_AVAILABILITYZONE_ID"
SetEnv "UKF_TEST_AVAILABILITYZONE_NAME"
SetEnv "UKF_TEST_AVAILABILITYZONE_NAME"
SetEnv "UKF_TEST_NETWORK_NAME"
SetEnv "UKF_TEST_DHCP_AVAILABILITYZONE_ID"
SetEnv "UKF_TEST_INSTANCE_NAME"
SetEnv "UKF_TEST_FLOATINGIP_ID"
# SetEnv "UKF_TEST_VPN_NAME"