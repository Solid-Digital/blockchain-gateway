const strength = {
  0: "Your password is not secure enough",
  1: "Your password is weak",
  2: "Still weak",
  3: "Pretty secure, but you could do better",
  4: "This is a really strong password!"
}

const email = document.getElementById('email');
const password = document.getElementById('password');
const passwordConfirmation = document.getElementById('password-confirmation');
const meter = document.getElementById('password-strength-meter');
const text = document.getElementById('password-strength-text');
const lengthValidationIcon = document.getElementById('length-validation__icon');
const upperCaseValidationIcon = document.getElementById('uppercase-validation__icon');
const specialCharacterValidationIcon = document.getElementById('special-char-validation__icon');
const submitButton = document.getElementById('submit-button');

if (submitButton !== null) {
  submitButton.disabled = true;
}

function performValidation() {
  const val = password.value;
  const confirmPasswordVal = passwordConfirmation.value;

  const specialXterRegex = /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;
  const upperCaseRegex = /[A-Z]/;
  const minLength = 8;
  const checkPassColour = "#6FCF97";
  const checkFailedColour = "#47515D";
  const result = zxcvbn(val);

  // Update the password strength meter
  meter.value = result.score;

  // Update the text indicator
  val !== "" ? text.innerHTML = strength[result.score] : submitButton.disabled = true;

  // check if password contains capital letter
  upperCaseRegex.test(val) ? upperCaseValidationIcon.style.color = checkPassColour : upperCaseValidationIcon.style.color = checkFailedColour;

  // check if password contains special character
  specialXterRegex.test(val) ? specialCharacterValidationIcon.style.color = checkPassColour : specialCharacterValidationIcon.style.color = checkFailedColour;

  // check if password is at least 8 characters
  val.length >= minLength ? lengthValidationIcon.style.color = checkPassColour : lengthValidationIcon.style.color = checkFailedColour;

  // disable submit button till all validations are complete
  (email === null || email.value !== "") && upperCaseRegex.test(val) && specialXterRegex.test(val) && val.length >= minLength & (val === confirmPasswordVal) ?
    submitButton.disabled = false : submitButton.disabled = true;
}

if (password !== null) {
  password.addEventListener('input', performValidation);
}

if (passwordConfirmation !== null) {
  passwordConfirmation.addEventListener('input', performValidation);
}

