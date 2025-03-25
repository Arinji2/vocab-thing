import {
  PersistedClient,
  Persister,
} from "@tanstack/react-query-persist-client";
import fs from "fs/promises";

const PERSIST_FILE_PATH = "./query-cache.json";

export function createFSPersister(): Persister {
  return {
    async persistClient(client: PersistedClient): Promise<void> {
      try {
        const data = JSON.stringify(client);
        await fs.writeFile(PERSIST_FILE_PATH, data, "utf-8");
      } catch (error) {
        console.error("Error persisting query client:", error);
      }
    },

    async restoreClient(): Promise<PersistedClient | undefined> {
      try {
        const data = await fs.readFile(PERSIST_FILE_PATH, "utf-8");
        return JSON.parse(data);
      } catch (error: any) {
        if (error.code !== "ENOENT") {
          console.error("Error restoring query client:", error);
        }
        return undefined;
      }
    },

    async removeClient(): Promise<void> {
      try {
        await fs.unlink(PERSIST_FILE_PATH);
      } catch (error: any) {
        if (error.code !== "ENOENT") {
          console.error("Error removing persisted query client:", error);
        }
      }
    },
  };
}
