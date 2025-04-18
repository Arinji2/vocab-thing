export function generateTitle(title?: string) {
  const defaultTitle = "Vocab Thing";
  if (!title) {
    return defaultTitle;
  }
  return `${title} | ${defaultTitle}`;
}
export function generateDescription(description?: string) {
  const defaultDescription =
    "Save words and phrases you find on the internet, and use them in the future effortlessly.";
  if (!description) {
    return defaultDescription;
  }
  return `${description}. ${defaultDescription}`;
}
