document.getElementById("viewTopicsBtn").addEventListener("click", async () => {
try {

	const response = await fetch('/private/alltopics');
	if (!response.ok) {
    	throw new Error(`HTTP error! Status: ${response.status}`);
}


const topics = await response.json();

console.log(topics); 

if (!Array.isArray(topics)) {
    throw new Error('Invalid response format');
}


const topicsContainer = document.getElementById("topics");
topicsContainer.innerHTML = "";


if (topics.length === 0) {
    topicsContainer.innerHTML = "<p>No topics available.</p>";
    return;
}


topics.forEach(topic => {
    const topicElement = document.createElement("div");
    topicElement.classList.add("topic");

    
    const topicName = topic.TopicName || "No Topic Name";
    const description = topic.description || "No Description";
    const content = topic.content || "No Content";

   
    
    topicElement.innerHTML = `
        <p><strong>${topicName}</strong> ${description}</p>
        `;
    topicsContainer.appendChild(topicElement);
});

} catch (error) {
console.error("Error loading topics:", error);
alert(`Failed to load topics: ${error.message}`);
}
});