// Simple test script to verify API connection
const testAPI = async () => {
  try {
    console.log("Testing API connection...");

    const response = await fetch("/api/shorten", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        url: "https://www.google.com",
      }),
    });

    console.log("Response status:", response.status);
    console.log("Response ok:", response.ok);

    if (response.ok) {
      const data = await response.json();
      console.log("Success! Response data:", data);
    } else {
      const errorText = await response.text();
      console.error("Error response:", errorText);
    }
  } catch (error) {
    console.error("Fetch failed:", error);
  }
};

// Run the test
testAPI();

export default testAPI;
