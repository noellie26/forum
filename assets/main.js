// Function to check if at least one of the tags is selected when creating a new post
function validateForm() {
    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    let checkedCount = 0;
    for (let i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            checkedCount++;
        }
    }
    if (checkedCount === 0) {
        alert("Please select at least one tag.");
        return false; // Prevent form submission
    }
    return true; // Allow form submission
}

// Function to convert ISO 8601 date-time format to a readable format
function convertDateTime(dateTimeString) {
    const dateTime = new Date(dateTimeString);
    const options = { day: 'numeric', month: 'long', year: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric' };
    return dateTime.toLocaleDateString('fr-FR', options);
}

// Function to update the date-time elements on page load
window.addEventListener('load', function () {
    const dateElements = document.querySelectorAll('.creationDate');
    dateElements.forEach(function (element) {
        const isoDateTime = element.textContent.trim();
        element.textContent = convertDateTime(isoDateTime);
    });
});

function toggleVisibility(id) {
    const registerForm = document.getElementById(id);
    if (registerForm.style.display !== 'block') {
        registerForm.style.display = 'block';
    } else {
        registerForm.style.display = 'none';
    }
}
function validatePassword() {
    var password = document.getElementById("password").value;
    var confirm_password = document.getElementById("confirm_password").value;
    var error_message = document.getElementById("password_error");

    if (password !== confirm_password) {
        error_message.style.display = "block";
        return false;
    } else {
        error_message.style.display = "none";
        return true;
    }
}

document.getElementById("inscription_form").addEventListener("submit", function(event) {
    if (!validatePassword()) {
        event.preventDefault();
    }
});
function togglePasswordVisibility() {
var passwordInput = document.getElementById("password");
var confirmInput = document.getElementById("confirm_password");
var showPasswordCheckbox = document.getElementById("show_password");

if (showPasswordCheckbox.checked) {
    passwordInput.type = "text";
    confirmInput.type = "text";
} else {
    passwordInput.type = "password";
    confirmInput.type = "password";
}
}


