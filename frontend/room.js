const copyButton = document.getElementById('copy-link');
const status = document.getElementById('status');
const chat = document.getElementById('chat');
const turn = document.getElementById('turn');
const errorMessage = document.getElementById('error');
const messageInput = document.getElementById('message');
const sendButton = document.getElementById('send');
const homeButton = document.getElementById('home')
const timerDisplay = document.getElementById('timer');
const modal = document.getElementsByClassName('modal')[0];

let timeLeft = 30;
let timer;
let ws;
homeButton.style.display = "none";
modal.style.display = "none";

function showToast() {
    let toast = document.getElementById("toast");
    toast.classList.add("show");

    setTimeout(() => {
        toast.classList.remove("show");
    }, 4500);
}


function startTimer() {
    clearInterval(timer);
    timeLeft = 30;
    updateTimerDisplay();

    timer = setInterval(() => {
        timeLeft--;
        updateTimerDisplay();

        if (timeLeft <= 0) {
            clearInterval(timer);
            ws.send(JSON.stringify({type: "timeout", text: "time is over"}));
        }
    }, 1000);
}

function updateTimerDisplay() {
    timerDisplay.textContent = `Время на ход: ${timeLeft} сек`;
}

copyButton.addEventListener('click', () => {
    navigator.clipboard.writeText(window.location.href);
    showToast();

});

const urlParams = new URLSearchParams(window.location.search);
const roomID = urlParams.get('roomID');

if (roomID) {
    const wsUrl = `ws://localhost:8080/ws/${roomID}`;
    ws = new WebSocket(wsUrl);

    ws.addEventListener('open', () => {
        status.textContent = 'Подключено, ждем второго игрока...';
    });

    ws.addEventListener('message', (event) => {
        errorMessage.textContent = ""
        let data;
        try {
            data = JSON.parse(event.data);
        } catch (error) {
            data = {"type": "message", "text": event.data};
        }
        if (data.connected !== undefined) {
            if (data.connected < 2) {
                status.textContent = `Ждём второго игрока... (${data.connected}/2)`;
                messageInput.disabled = true;
                sendButton.disabled = true;
            } else {
                ws.send(JSON.stringify({type: "turn", text: "get turn"}));
                status.textContent = "Два игрока в комнате! Можно начинать.";
                modal.style.display = "";
                copyButton.style.display = "none";
                status.style.display = "none";
                document.body.style.backgroundColor = "#DCDCDC";
                messageInput.disabled = false;
                sendButton.disabled = false;
                startTimer();
            }

        } else if (data.type !== undefined) {
            if (data.type === "result") {
                messageInput.style.display = "none";
                sendButton.style.display = "none";
                modal.style.display = "none";
                errorMessage.style.display = "none";
                turn.style.display = "none"
                if (data.text === "Вы проиграли") {
                    chat.textContent = "Вы проиграли :( Вы не успели!";
                    document.body.style.removeProperty('background');
                    document.body.style.background = 'linear-gradient(270deg, #2c0101, #440a0a, #501111)';
                    document.body.style.backgroundSize = '200% 200%';
                    document.body.style.animation = 'gradientAnimation 11s ease infinite';
                } else if (data.text === "Вы победили") {
                    chat.textContent = "Вы победили :) Соперник не успел!";
                    document.body.style.removeProperty('background');
                    document.body.style.background = 'linear-gradient(270deg, #071e02, #052805, #0b2f04)';
                    document.body.style.backgroundSize = '200% 200%';
                    document.body.style.animation = 'gradientAnimation 11s ease infinite';
                }
                homeButton.style.display = "";
                homeButton.textContent = "На главную";

            } else if (data.type === "message") {
                clearInterval(timer);
                let specialLetters = ["ь", "ъ", "ы"];
                let chars = data.text.split("");
                let index = chars.length - 1;
                if (specialLetters.includes(chars[index].toLowerCase()) && index > 0) {
                    index -= 1;
                }
                chars[index] = `${chars[index].toUpperCase()}`;
                chat.textContent = "Последние слово: " + chars.join("");
                startTimer();
            } else if (data.type === "error") {
                if (data.text === "Комната уже заполнена!") {
                    timerDisplay.style.display = "none"
                    status.textContent = "Комната уже заполнена!";
                } else {
                    errorMessage.textContent = data.text;
                }
            } else if (data.type === "turn") {
                turn.textContent = data.text;
            }
        }
    });

    sendButton.addEventListener('click', () => {
        const message = messageInput.value.trim();
        if (message) {
            ws.send(JSON.stringify({type: "message", text: message.toLowerCase().trim()}));
            messageInput.value = '';
        }
    });

    messageInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter' && !messageInput.disabled) {
            sendButton.click();
        }
    });
}
