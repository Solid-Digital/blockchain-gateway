var strength = {
  0: "Your password is not secure enough",
  1: "Your password is weak",
  2: "Still weak",
  3: "Pretty secure, but you could do better",
  4: "This is a really strong password!"
}

const password = document.getElementById('password');
const meter = document.getElementById('password-strength-meter');
const text = document.getElementById('password-strength-text');
const lengthValidationIcon = document.getElementById('length-validation__icon');
const upperCaseValidationIcon = document.getElementById('uppercase-validation__icon');
const specialCharacterValidationIcon = document.getElementById('special-char-validation__icon');
const passwordCheckerDiv = document.getElementById('password-checker');
const submitButton = document.getElementById('submit-button');

if (submitButton !== null) {
  submitButton.disabled = true;
}

if (password !== null ) {
  password.addEventListener('input', function () {
    var val = password.value;
    var result = zxcvbn(val);

    // Update the password strength meter
    meter.value = result.score;

    // Update the text indicator
    val !== "" ? text.innerHTML = strength[result.score] : submitButton.disabled = true;

    // check if password contains capital letter
    /[A-Z]/.test(val) ? upperCaseValidationIcon.style.color = "#6FCF97" : upperCaseValidationIcon.style.color = "#47515D";

    // check if password contains special character
    /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(val) ? specialCharacterValidationIcon.style.color = "#6FCF97" : specialCharacterValidationIcon.style.color = "#47515D";

    // check if password is at least 12 characters
    val.length >= 12 ? lengthValidationIcon.style.color = "#6FCF97" : lengthValidationIcon.style.color = "#47515D";

    // disable submit button till all validations are complete
    /[A-Z]/.test(val) && /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(val) && val.length >= 12 ? submitButton.disabled = false : submitButton.disabled = true;
  });
}