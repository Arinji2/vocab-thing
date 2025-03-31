export function getApiURL() {
  const envURL = process.env.API_URL;
  if (envURL) {
    return envURL;
  } else {
    console.error("API_URL is not set");
    return "http://localhost:8080";
  }
}
