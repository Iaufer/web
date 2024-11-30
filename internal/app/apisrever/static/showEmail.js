console.log('Script loaded');

async function fetchUserName() {
    try {
        const response = await fetch('/private/whoami', { credentials: 'include' });
        if (!response.ok) throw new Error('Failed to fetch user data');
        const data = await response.json();
        document.getElementById('username').innerHTML = `
            <span>${data.email}</span>
            <a href="/own-topics.html">OWN topics</a>
        `;
    } catch (err) {
        console.error('Error fetching user data:', err);
    }
}

// Загрузка данных при открытии страницы
fetchUserName();

// Управление модальным окном
const newTopicModal = document.getElementById("newTopicModal");
const newTopicBtn = document.getElementById("newTopicBtn");
const closeModalBtn = document.getElementById("closeModalBtn");

newTopicBtn.onclick = function() {
    newTopicModal.style.display = "block";
}

closeModalBtn.onclick = function() {
    newTopicModal.style.display = "none";
}

window.onclick = function(event) {
    if (event.target === newTopicModal) {
        newTopicModal.style.display = "none";
    }
}