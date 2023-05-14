const urlParams = new URLSearchParams(window.location.search);
const userID = urlParams.get('userID');

const socket = new WebSocket(`ws://localhost:8081/ws?userID=${userID}`);

socket.addEventListener('open', function(event) {
  console.log('WebSocket connection established.');
});

socket.addEventListener('message', function(event) {
  console.log('Received message: ' + event.data);
});

socket.addEventListener('close', function(event) {
  console.log('WebSocket connection closed.');
});


// Send a message to the server
function sendMessage(message) {
  socket.send(message);
}
