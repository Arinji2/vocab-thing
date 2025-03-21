import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <div className="p-2">
      <h3 className="text-white">Welcome Home!!!</h3>
      <h3 className="text-white font-medium">This is a medium font</h3>
    </div>
  );
}
