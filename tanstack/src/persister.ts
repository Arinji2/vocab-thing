import { PersistedClient } from "@tanstack/react-query-persist-client";
import path from "path";
import fs from "fs/promises";

type JsonPersisterOptions = {
  filePath?: string;
  filename?: string;
};
export function createJsonFilePersister({
  filePath = process.cwd(),
  filename = "query-cache.json",
}: JsonPersisterOptions = {}) {
  const fullPath = path.resolve(filePath, filename);

  return {
    persistClient: async (client: PersistedClient) => {
      try {
        await fs.mkdir(path.dirname(fullPath), { recursive: true });

        await fs.writeFile(fullPath, JSON.stringify(client, null, 2), "utf8");
      } catch (error) {
        console.error("Error persisting client to JSON:", error);
      }
    },

    restoreClient: async () => {
      try {
        try {
          await fs.access(fullPath);
        } catch {
          return undefined;
        }

        const fileContents = await fs.readFile(fullPath, "utf8");
        return fileContents ? JSON.parse(fileContents) : undefined;
      } catch (error) {
        console.error("Error restoring client from JSON:", error);
        return undefined;
      }
    },

    removeClient: async () => {
      try {
        await fs.unlink(fullPath);
      } catch (error) {
        if ((error as NodeJS.ErrnoException).code !== "ENOENT") {
          console.error("Error removing client JSON:", error);
        }
      }
    },
  };
}
