export const seo = ({
  title,
  description,
  keywords,
  image,
  theme,
}: {
  title: string;
  description?: string;
  image?: string;
  keywords?: string;
  theme?: string;
}) => {
  const tags = [
    { title },
    { name: "description", content: description },
    { name: "keywords", content: keywords },
    { name: "theme-color", content: theme },
    { name: "twitter:title", content: title },
    { name: "twitter:description", content: description },
    { name: "twitter:creator", content: "@Arinji_i" },
    { name: "twitter:site", content: "@Arinji_i" },
    { name: "og:type", content: "website" },
    { name: "og:title", content: title },
    { name: "og:description", content: description },
    ...(image
      ? [
          { name: "twitter:image", content: image },
          { name: "twitter:card", content: "summary_large_image" },
          { name: "og:image", content: image },
        ]
      : []),
  ];

  return tags;
};
