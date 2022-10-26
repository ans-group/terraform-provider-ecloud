Param
(
    $Force = $false
)

$env:TF_ACC = "1"
function SetEnv($Name)
{
    if (((Test-Path env:\$Name) -ne $true) -or $Force)
    {
        New-Item -Name $Name -value (Read-Host -Prompt "Enter $Name") -ItemType Variable -Path Env: | Out-Null
    }
}

SetEnv "ANS_API_KEY"
SetEnv "ANS_TEST_VPC_REGION_ID"