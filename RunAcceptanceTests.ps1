Param
(
    [array]$Tests = @()
)
. ./SetTestEnv.ps1
$AdditionalArguments = @()
if ($Tests.Count -gt 0)
{
    $AdditionalArguments += "-run=" + ($Tests -join ",")
}

$env:TF_DEBUG = 1
go test -v -timeout=120m ./ecloud ($AdditionalArguments -join " ")