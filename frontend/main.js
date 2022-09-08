var alarmString = null;

// Select DOM element with create-alarm id
const createAlarm = document.querySelector(".create-alarm");

const turnOn = async () => {
  const response = await fetch('127.0.0.1:8031/api/on', {
    method: 'PATCH',
  });
}

const turnOff = async () => {
  const response = await fetch('127.0.0.1:8031/api/off', {
    method: 'PATCH',
  });
}

// Select DOM element of active alarm container
const activeAlarm = document.getElementById("active-alarm");
const clearAlarm = document.getElementById("clear-alarm");

// Select DOM element of active alarm text
const alarmTextContainer = document.getElementById("alarm-text");

const alarmText = (time) => `Alarm set at time ${time}`;

// Handle Create Alarm submit
const handleSubmit = (event) => {
  // Prevent default action of reloading the page
  event.preventDefault();
  const { hour, min } = document.forms[0];
  alarmString = getTimeString({
    hours: hour.value,
    minutes: min.value,
  });
  // Reset form after submit
  document.forms[0].reset();
  // Hide create alarm
  createAlarm.style.display = "none";
  // show active alarm with text
  activeAlarm.style.display = "block";
  alarmTextContainer.innerHTML = alarmText(alarmString);
};

const handleClear = () => {
  alarmString = "";
  activeAlarm.style.display = "none";
  createAlarm.style.display = "block";
  deleteAlarm();
};

const deleteAlarm = async () => {
  const response = await fetch('http://localhost:8031/api/clear', {
    method: 'DELETE',
  });
}

// Trigger handleClear on button click
clearAlarm.addEventListener("click", handleClear);
// Attach submit event to the form
document.forms[0].addEventListener("submit", handleSubmit);

// Function to convert time to string value
const getTimeString = ({ hours, minutes }) => {
  if (minutes / 10 < 1) {
    minutes = "0" + minutes;
  }
  return `${hours}:${minutes}`;
};

// Function to display current time on screen
const renderTime = () => {
  var currentTime = document.getElementById("current-time");
  const currentDate = new Date();
  var hours = currentDate.getHours();
  var minutes = currentDate.getMinutes();
  const timeString = getTimeString({ hours, minutes });
  currentTime.innerHTML = timeString;
};

// Update time every second
setInterval(renderTime, 1000);