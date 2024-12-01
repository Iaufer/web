document.getElementById("viewTopicsBtn").addEventListener("click", async () => {
    const topicsPerPage = 10; // Количество тем на странице
    let currentPage = 1; // Текущая страница
    let topics = []; // Список всех тем

    // Функция для получения всех тем
    async function fetchTopics() {
        try {
            const response = await fetch('/private/alltopics');
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            const data = await response.json();
            if (!Array.isArray(data)) {
                throw new Error('Invalid response format');
            }
            return data;
        } catch (error) {
            console.error("Error loading topics:", error);
            alert(`Failed to load topics: ${error.message}`);
            return [];
        }
    }

    // Функция для отображения тем на текущей странице
    function renderTopics() {
        const topicsContainer = document.getElementById("topics");
        topicsContainer.innerHTML = "";

        const startIndex = (currentPage - 1) * topicsPerPage;
        const endIndex = Math.min(startIndex + topicsPerPage, topics.length);

        if (topics.length === 0) {
            topicsContainer.innerHTML = "<p>No topics available.</p>";
            return;
        }

        for (let i = startIndex; i < endIndex; i++) {
            const topic = topics[i];
            const topicElement = document.createElement("div");
            topicElement.classList.add("topic");

            const topicName = topic.TopicName || "No Topic Name";
            const description = topic.description || "No Description";

            topicElement.innerHTML = `
                <p>
                    <strong>
                        <a href="/private/topic/${topic.id}" style="text-decoration: none; color: inherit;">
                            ${topicName}
                        </a>
                    </strong>
                    ${description}
                </p>
            `;
            topicsContainer.appendChild(topicElement);
        }
    }

    // Функция для отображения кнопок пагинации
    function renderPagination() {
        const paginationContainer = document.getElementById("pagination");
        paginationContainer.innerHTML = ""; // Очистка предыдущих кнопок

        const totalPages = Math.ceil(topics.length / topicsPerPage);

        for (let i = 1; i <= totalPages; i++) {
            const pageButton = document.createElement("button");
            pageButton.textContent = i;
            pageButton.classList.add("page-btn");
            if (i === currentPage) {
                pageButton.classList.add("active");
            }

            pageButton.addEventListener("click", () => {
                currentPage = i;
                renderTopics();
                renderPagination();
            });

            paginationContainer.appendChild(pageButton);
        }
    }


    // Загрузка тем и рендеринг страницы
    topics = await fetchTopics();
    if (topics.length > 0) {
        renderTopics();
        renderPagination();
    }
});