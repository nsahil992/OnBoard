// Load employees on page load
document.addEventListener('DOMContentLoaded', function() {
    loadEmployees();
});

// Handle form submission
document.getElementById('employeeForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const formData = new FormData(e.target);
    const employeeData = {
        first_name: formData.get('firstName'),
        last_name: formData.get('lastName'),
        email: formData.get('email'),
        phone: formData.get('phone'),
        department: formData.get('department'),
        job_title: formData.get('jobTitle'),
        hire_date: formData.get('hireDate'),
        salary: parseFloat(formData.get('salary'))
    };

    try {
        const response = await fetch('/api/employees', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(employeeData)
        });

        const messageDiv = document.getElementById('message');

        if (response.ok) {
            messageDiv.innerHTML = '<div class="success-message">Employee added successfully!</div>';
            e.target.reset();
            loadEmployees(); // Refresh the employee list
        } else {
            const error = await response.text();
            messageDiv.innerHTML = `<div class="error-message">Error: ${error}</div>`;
        }
    } catch (error) {
        document.getElementById('message').innerHTML = `<div class="error-message">Error: ${error.message}</div>`;
    }
});

// Load and display employees
async function loadEmployees() {
    try {
        const response = await fetch('/api/employees');
        const employees = await response.json();

        const employeesList = document.getElementById('employeesList');
        const employeeCount = document.getElementById('employeeCount');

        employeeCount.textContent = employees.length;

        if (employees.length === 0) {
            employeesList.innerHTML = '<div class="loading">No employees found.</div>';
            return;
        }

        const employeesHtml = employees.map(emp => `
            <div class="employee-card">
                <div class="employee-card-header">
                    <div class="employee-name">${emp.first_name} ${emp.last_name}</div>
                    <button class="delete-btn" onclick="deleteEmployee(${emp.id})">Delete</button>
                </div>
                <div class="employee-info">
                    <div class="info-item">
                        <span class="info-label">Email:</span>
                        <span>${emp.email}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">Phone:</span>
                        <span>${emp.phone}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">Department:</span>
                        <span>${emp.department}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">Job Title:</span>
                        <span>${emp.job_title}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">Hire Date:</span>
                        <span>${new Date(emp.hire_date).toLocaleDateString()}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">Salary:</span>
                        <span>$${parseFloat(emp.salary).toLocaleString()}</span>
                    </div>
                </div>
            </div>
        `).join('');

        employeesList.innerHTML = `<div class="employee-grid">${employeesHtml}</div>`;
    } catch (error) {
        document.getElementById('employeesList').innerHTML = `<div class="error-message">Error loading employees: ${error.message}</div>`;
    }
}

// Delete employee
async function deleteEmployee(id) {
    if (!confirm('Are you sure you want to delete this employee?')) return;

    try {
        const response = await fetch(`/api/employees/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showMessage('Employee deleted successfully!', 'success');
            loadEmployees();
        } else {
            const error = await response.text();
            showMessage(`Error: ${error}`, 'error');
        }
    } catch (error) {
        showMessage(`Error: ${error.message}`, 'error');
    }
}

function showMessage(text, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.innerHTML = `
        <div class="${type === 'success' ? 'success-message' : 'error-message'}">
            ${text}
        </div>
    `;
}
