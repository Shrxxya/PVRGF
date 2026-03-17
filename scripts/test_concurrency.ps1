$baseUrl = "http://localhost:8080/api"
$userId = 1

Write-Host "Starting Parallel Request Simulation..." -ForegroundColor Cyan
Write-Host "Simulating 5 simultaneous users adding passwords (each has 500ms delay)..."

$passwords = @(
    @{ label = "gmail"; password = "Password123!" },
    @{ label = "facebook"; password = "Pass!234Word" },
    @{ label = "instagram"; password = "Insta!Gram2024" },
    @{ label = "linkedin"; password = "Linked!In2024" },
    @{ label = "github"; password = "Git!Hub2024" }
)

$startTime = Get-Date

# Start 5 parallel jobs
$jobs = $passwords | ForEach-Object {
    $payload = @{
        userId = $userId
        label = $_.label
        password = $_.password
    } | ConvertTo-Json

    Start-Job -ScriptBlock {
        param($url, $data)
        Invoke-RestMethod -Uri $url -Method Post -Body $data -ContentType "application/json"
    } -ArgumentList "$baseUrl/passwords", $payload
}

# Wait for all jobs to finish
Write-Host "Waiting for workers to complete..."
$results = $jobs | Wait-Job | Receive-Job
Remove-Job $jobs

$endTime = Get-Date
$duration = ($endTime - $startTime).TotalSeconds

Write-Host "`nSimulation Complete!" -ForegroundColor Green
Write-Host "Total Time Taken: $duration seconds" -ForegroundColor Yellow
Write-Host "Sequential Expected Time: 2.5 seconds"
Write-Host "Worker Pool Efficiency: $([math]::Round(2.5 - $duration, 2)) seconds saved."

Write-Host "`nCheck the server terminal logs to see workers (0-4) processing these in parallel!"
