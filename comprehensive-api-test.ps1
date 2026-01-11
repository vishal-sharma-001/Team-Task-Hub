#!/usr/bin/env pwsh
# Comprehensive API Test Suite

$BASE_URL = "http://localhost:8080/api"
$RESULTS = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [hashtable]$Headers = @{}
    )
    
    try {
        $params = @{
            Uri = "$BASE_URL$Endpoint"
            Method = $Method
            UseBasicParsing = $true
            TimeoutSec = 10
        }
        
        if ($Body) {
            $params['Body'] = $Body | ConvertTo-Json -Depth 10
            $params['ContentType'] = 'application/json'
        }
        
        if ($Headers.Count -gt 0) {
            $params['Headers'] = $Headers
        }
        
        $response = Invoke-WebRequest @params
        $data = $response.Content | ConvertFrom-Json
        
        $result = @{
            Test = $Name
            Status = $response.StatusCode
            Success = $true
            Response = $data
        }
    } catch {
        $result = @{
            Test = $Name
            Status = "ERROR"
            Success = $false
            Error = $_.Exception.Message
        }
    }
    
    return $result
}

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "COMPREHENSIVE API TEST SUITE" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# ============ PART 1: AUTH TESTS ============
Write-Host "PART 1: AUTHENTICATION TESTS" -ForegroundColor Yellow
Write-Host "================================" -ForegroundColor Yellow

# Prepare unique test identity
$TIMESTAMP = [DateTime]::UtcNow.ToString("yyyyMMddHHmmss")
$EMAIL = "testuser_$TIMESTAMP@example.com"

# Signup
Write-Host "`n[1] Testing POST /api/auth/signup..."
$signupData = @{
    email = $EMAIL
    password = "TestPassword123!"
}
$signup = Test-Endpoint -Name "Signup" -Method "POST" -Endpoint "/auth/signup" -Body $signupData
Write-Host "Status: $($signup.Status) | Success: $($signup.Success)"
if ($signup.Success) {
    Write-Host "Response: $($signup.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
    $TOKEN = $signup.Response.data.token
    $USER = $signup.Response.data.user
} else {
    Write-Host "Error: $($signup.Error)" -ForegroundColor Red
}
Write-Host ""

# Login
Write-Host "[2] Testing POST /api/auth/login..."
$loginData = @{
    email = $EMAIL
    password = "TestPassword123!"
}
$login = Test-Endpoint -Name "Login" -Method "POST" -Endpoint "/auth/login" -Body $loginData
Write-Host "Status: $($login.Status) | Success: $($login.Success)"
if ($login.Success) {
    Write-Host "Response: $($login.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
    $TOKEN = $login.Response.data.token
} else {
    Write-Host "Error: $($login.Error)" -ForegroundColor Red
}
Write-Host ""

# Setup auth headers for protected routes
$authHeaders = @{ "Authorization" = "Bearer $TOKEN" }
if (-not $TOKEN -or $TOKEN -eq "") {
    Write-Host "Auth token missing; stopping protected tests." -ForegroundColor Red
    exit 1
}

# Get Profile
Write-Host "[3] Testing GET /api/auth/me..."
$profile = Test-Endpoint -Name "GetProfile" -Method "GET" -Endpoint "/auth/me" -Headers $authHeaders
Write-Host "Status: $($profile.Status) | Success: $($profile.Success)"
if ($profile.Success) {
    Write-Host "Profile: $($profile.Response.data.email)" -ForegroundColor Green
    $USER = $profile.Response.data
} else {
    Write-Host "Error: $($profile.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Profile
Write-Host "[4] Testing PUT /api/auth/me..."
$updateData = @{
    first_name = "Test"
    last_name = "User"
}
$updateProfile = Test-Endpoint -Name "UpdateProfile" -Method "PUT" -Endpoint "/auth/me" -Body $updateData -Headers $authHeaders
Write-Host "Status: $($updateProfile.Status) | Success: $($updateProfile.Success)"
if ($updateProfile.Success) {
    Write-Host "Updated: $($updateProfile.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateProfile.Error)" -ForegroundColor Red
}
Write-Host ""

# List Users
Write-Host "[5] Testing GET /api/users..."
$users = Test-Endpoint -Name "ListUsers" -Method "GET" -Endpoint "/users" -Headers $authHeaders
Write-Host "Status: $($users.Status) | Success: $($users.Success)"
if ($users.Success) {
    Write-Host "Users found: $(@($users.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($users.Error)" -ForegroundColor Red
}
Write-Host ""

# ============ PART 2: PROJECT TESTS ============
Write-Host "PART 2: PROJECT TESTS" -ForegroundColor Yellow
Write-Host "======================" -ForegroundColor Yellow

# Create Project
Write-Host "`n[6] Testing POST /api/projects..."
$projectData = @{
    name = "Test Project"
    description = "A test project for API validation"
}
$createProject = Test-Endpoint -Name "CreateProject" -Method "POST" -Endpoint "/projects" -Body $projectData -Headers $authHeaders
Write-Host "Status: $($createProject.Status) | Success: $($createProject.Success)"
if ($createProject.Success) {
    $PROJECT_ID = $createProject.Response.data.id
    Write-Host "Project Created: $PROJECT_ID" -ForegroundColor Green
    Write-Host "Response: $($createProject.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} else {
    Write-Host "Error: $($createProject.Error)" -ForegroundColor Red
}
Write-Host ""

# List Projects
Write-Host "[7] Testing GET /api/projects..."
$listProjects = Test-Endpoint -Name "ListProjects" -Method "GET" -Endpoint "/projects" -Headers $authHeaders
Write-Host "Status: $($listProjects.Status) | Success: $($listProjects.Success)"
if ($listProjects.Success) {
    Write-Host "Projects found: $(@($listProjects.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($listProjects.Error)" -ForegroundColor Red
}
Write-Host ""

# Get Project
Write-Host "[8] Testing GET /api/projects/{project_id}..."
$getProject = Test-Endpoint -Name "GetProject" -Method "GET" -Endpoint "/projects/$PROJECT_ID" -Headers $authHeaders
Write-Host "Status: $($getProject.Status) | Success: $($getProject.Success)"
if ($getProject.Success) {
    Write-Host "Project: $($getProject.Response.data.name)" -ForegroundColor Green
    Write-Host "Creator ID: $($getProject.Response.data.created_by_id)" -ForegroundColor Green
    Write-Host "Full Response: $($getProject.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} else {
    Write-Host "Error: $($getProject.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Project
Write-Host "[9] Testing PUT /api/projects/{project_id}..."
$updateProjectData = @{
    name = "Updated Test Project"
    description = "Updated description"
}
$updateProject = Test-Endpoint -Name "UpdateProject" -Method "PUT" -Endpoint "/projects/$PROJECT_ID" -Body $updateProjectData -Headers $authHeaders
Write-Host "Status: $($updateProject.Status) | Success: $($updateProject.Success)"
if ($updateProject.Success) {
    Write-Host "Updated: $($updateProject.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateProject.Error)" -ForegroundColor Red
}
Write-Host ""

# ============ PART 3: TASK TESTS ============
Write-Host "PART 3: TASK TESTS" -ForegroundColor Yellow
Write-Host "==================" -ForegroundColor Yellow

# Create Task
Write-Host "`n[10] Testing POST /api/projects/{project_id}/tasks..."
$taskData = @{
    title = "Test Task"
    description = "A test task"
    priority = "MEDIUM"
    assignee_id = $USER.id
}
$createTask = Test-Endpoint -Name "CreateTask" -Method "POST" -Endpoint "/projects/$PROJECT_ID/tasks" -Body $taskData -Headers $authHeaders
Write-Host "Status: $($createTask.Status) | Success: $($createTask.Success)"
if ($createTask.Success) {
    $TASK_ID = $createTask.Response.data.id
    Write-Host "Task Created: $TASK_ID" -ForegroundColor Green
    Write-Host "Response: $($createTask.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} else {
    Write-Host "Error: $($createTask.Error)" -ForegroundColor Red
}
Write-Host ""

# List Project Tasks
Write-Host "[11] Testing GET /api/projects/{project_id}/tasks..."
$listTasks = Test-Endpoint -Name "ListTasks" -Method "GET" -Endpoint "/projects/$PROJECT_ID/tasks" -Headers $authHeaders
Write-Host "Status: $($listTasks.Status) | Success: $($listTasks.Success)"
if ($listTasks.Success) {
    Write-Host "Tasks found: $(@($listTasks.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($listTasks.Error)" -ForegroundColor Red
}
Write-Host ""

# Get Assigned Tasks
Write-Host "[12] Testing GET /api/tasks/assigned..."
$assignedTasks = Test-Endpoint -Name "ListAssignedTasks" -Method "GET" -Endpoint "/tasks/assigned" -Headers $authHeaders
Write-Host "Status: $($assignedTasks.Status) | Success: $($assignedTasks.Success)"
if ($assignedTasks.Success) {
    Write-Host "Assigned tasks: $(@($assignedTasks.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($assignedTasks.Error)" -ForegroundColor Red
}
Write-Host ""

# Get Task
Write-Host "[13] Testing GET /api/tasks/{task_id}..."
$getTask = Test-Endpoint -Name "GetTask" -Method "GET" -Endpoint "/tasks/$TASK_ID" -Headers $authHeaders
Write-Host "Status: $($getTask.Status) | Success: $($getTask.Success)"
if ($getTask.Success) {
    Write-Host "Task: $($getTask.Response.data.title)" -ForegroundColor Green
    Write-Host "Created By: $($getTask.Response.data.created_by_id)" -ForegroundColor Green
    Write-Host "Assigned To: $($getTask.Response.data.assignee_id)" -ForegroundColor Green
    Write-Host "Full Response: $($getTask.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} else {
    Write-Host "Error: $($getTask.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Task
Write-Host "[14] Testing PUT /api/tasks/{task_id}..."
$updateTaskData = @{
    title = "Updated Task Title"
    description = "Updated description"
}
$updateTask = Test-Endpoint -Name "UpdateTask" -Method "PUT" -Endpoint "/tasks/$TASK_ID" -Body $updateTaskData -Headers $authHeaders
Write-Host "Status: $($updateTask.Status) | Success: $($updateTask.Success)"
if ($updateTask.Success) {
    Write-Host "Updated: $($updateTask.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateTask.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Task Status
Write-Host "[15] Testing PATCH /api/tasks/{task_id}/status..."
$statusData = @{ status = "IN_PROGRESS" }
$updateStatus = Test-Endpoint -Name "UpdateTaskStatus" -Method "PATCH" -Endpoint "/tasks/$TASK_ID/status" -Body $statusData -Headers $authHeaders
Write-Host "Status: $($updateStatus.Status) | Success: $($updateStatus.Success)"
if ($updateStatus.Success) {
    Write-Host "Status updated: $($updateStatus.Response.data.status)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateStatus.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Task Priority
Write-Host "[16] Testing PATCH /api/tasks/{task_id}/priority..."
$priorityData = @{ priority = "HIGH" }
$updatePriority = Test-Endpoint -Name "UpdateTaskPriority" -Method "PATCH" -Endpoint "/tasks/$TASK_ID/priority" -Body $priorityData -Headers $authHeaders
Write-Host "Status: $($updatePriority.Status) | Success: $($updatePriority.Success)"
if ($updatePriority.Success) {
    Write-Host "Priority updated: $($updatePriority.Response.data.priority)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updatePriority.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Task Assignee
Write-Host "[17] Testing PATCH /api/tasks/{task_id}/assignee..."
$assigneeData = @{ assignee_id = $USER.id }
$updateAssignee = Test-Endpoint -Name "UpdateTaskAssignee" -Method "PATCH" -Endpoint "/tasks/$TASK_ID/assignee" -Body $assigneeData -Headers $authHeaders
Write-Host "Status: $($updateAssignee.Status) | Success: $($updateAssignee.Success)"
if ($updateAssignee.Success) {
    Write-Host "Assignee updated: $($updateAssignee.Response.data.assigned_by_id)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateAssignee.Error)" -ForegroundColor Red
}
Write-Host ""

# ============ PART 4: COMMENT TESTS ============
Write-Host "PART 4: COMMENT TESTS" -ForegroundColor Yellow
Write-Host "=====================" -ForegroundColor Yellow

# Create Comment
Write-Host "`n[18] Testing POST /api/tasks/{task_id}/comments..."
$commentData = @{
    content = "This is a test comment"
}
$createComment = Test-Endpoint -Name "CreateComment" -Method "POST" -Endpoint "/tasks/$TASK_ID/comments" -Body $commentData -Headers $authHeaders
Write-Host "Status: $($createComment.Status) | Success: $($createComment.Success)"
if ($createComment.Success) {
    $COMMENT_ID = $createComment.Response.data.id
    Write-Host "Comment Created: $COMMENT_ID" -ForegroundColor Green
    Write-Host "Response: $($createComment.Response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} else {
    Write-Host "Error: $($createComment.Error)" -ForegroundColor Red
}
Write-Host ""

# List Comments
Write-Host "[19] Testing GET /api/tasks/{task_id}/comments..."
$listComments = Test-Endpoint -Name "ListComments" -Method "GET" -Endpoint "/tasks/$TASK_ID/comments" -Headers $authHeaders
Write-Host "Status: $($listComments.Status) | Success: $($listComments.Success)"
if ($listComments.Success) {
    Write-Host "Comments found: $(@($listComments.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($listComments.Error)" -ForegroundColor Red
}
Write-Host ""

# Get Recent Comments
Write-Host "[20] Testing GET /api/comments/recent..."
$recentComments = Test-Endpoint -Name "ListRecentComments" -Method "GET" -Endpoint "/comments/recent" -Headers $authHeaders
Write-Host "Status: $($recentComments.Status) | Success: $($recentComments.Success)"
if ($recentComments.Success) {
    Write-Host "Recent comments: $(@($recentComments.Response.data).Count)" -ForegroundColor Green
} else {
    Write-Host "Error: $($recentComments.Error)" -ForegroundColor Red
}
Write-Host ""

# Update Comment
Write-Host "[21] Testing PUT /api/comments/{comment_id}..."
$updateCommentData = @{
    content = "Updated comment text"
}
$updateComment = Test-Endpoint -Name "UpdateComment" -Method "PUT" -Endpoint "/comments/$COMMENT_ID" -Body $updateCommentData -Headers $authHeaders
Write-Host "Status: $($updateComment.Status) | Success: $($updateComment.Success)"
if ($updateComment.Success) {
    Write-Host "Updated: $($updateComment.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($updateComment.Error)" -ForegroundColor Red
}
Write-Host ""

# Delete Comment
Write-Host "[22] Testing DELETE /api/comments/{comment_id}..."
$deleteComment = Test-Endpoint -Name "DeleteComment" -Method "DELETE" -Endpoint "/comments/$COMMENT_ID" -Headers $authHeaders
Write-Host "Status: $($deleteComment.Status) | Success: $($deleteComment.Success)"
if ($deleteComment.Success) {
    Write-Host "Deleted: $($deleteComment.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($deleteComment.Error)" -ForegroundColor Red
}
Write-Host ""

# ============ PART 5: DELETE TESTS ============
Write-Host "PART 5: CLEANUP/DELETE TESTS" -ForegroundColor Yellow
Write-Host "=============================" -ForegroundColor Yellow

# Delete Task
Write-Host "`n[23] Testing DELETE /api/tasks/{task_id}..."
$deleteTask = Test-Endpoint -Name "DeleteTask" -Method "DELETE" -Endpoint "/tasks/$TASK_ID" -Headers $authHeaders
Write-Host "Status: $($deleteTask.Status) | Success: $($deleteTask.Success)"
if ($deleteTask.Success) {
    Write-Host "Deleted: $($deleteTask.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($deleteTask.Error)" -ForegroundColor Red
}
Write-Host ""

# Delete Project
Write-Host "[24] Testing DELETE /api/projects/{project_id}..."
$deleteProject = Test-Endpoint -Name "DeleteProject" -Method "DELETE" -Endpoint "/projects/$PROJECT_ID" -Headers $authHeaders
Write-Host "Status: $($deleteProject.Status) | Success: $($deleteProject.Success)"
if ($deleteProject.Success) {
    Write-Host "Deleted: $($deleteProject.Response.message)" -ForegroundColor Green
} else {
    Write-Host "Error: $($deleteProject.Error)" -ForegroundColor Red
}
Write-Host ""

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "TEST COMPLETE" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
