console.log('Script loaded');

async function fetchUserName() {
    try {
        const response = await fetch('/private/premiumcontent', { credentials: 'include' });
        if (!response.ok) throw new Error('Failed to premium user data');
        const data = await response.json();
        document.getElementById('premium').innerHTML = `
            <span>Premium ✅</span>
        `;
    } catch (err) {
        console.error('Error fetching user data:', err);
    }
}

// Загрузка данных при открытии страницы
fetchUserName();
