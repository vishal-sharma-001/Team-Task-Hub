$headers = @{
    "Authorization" = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJlbWFpbCI6InNoYXJtYXZzMjAwMUBnbWFpbC5jb20iLCJleHAiOjE3NjgxMTQxNTMsImlhdCI6MTc2ODAyNzc1M30.oG92-xHJf12_UsDqz0naijDZA3AZSFW4C8ud9LCbCgs"
    "Content-Type" = "application/json"
}

$body = '{"status":"IN_PROGRESS"}'

try {
    Write-Host "Testing partial update: PUT /api/tasks/4 with {status:IN_PROGRESS}"
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/tasks/4" -Method PUT -Headers $headers -Body $body
    Write-Host "SUCCESS: Status code $($response.StatusCode)"
    Write-Host "Response:"
    $response.Content | ConvertFrom-Json | ConvertTo-Json
} catch {
    Write-Host "ERROR: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorContent = $reader.ReadToEnd()
        Write-Host "Response: $errorContent"
    }
}
