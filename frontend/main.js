const createButton = document.getElementById('create-room');

createButton.addEventListener('click', async () => {
    const response = await fetch('http://localhost:8080/create');
    const link = await response.text();

    const roomID = link.split('/').pop();
    window.location.href = `room.html?roomID=${roomID}`;
});
