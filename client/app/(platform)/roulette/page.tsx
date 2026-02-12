import { RouletteTableList } from "@/components/roulette/roulette-table-list";
import { PageHeader } from "@/components/shared/page-header";

export default function RoulettePage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Roulette"
        description="Choose a table to start playing."
      />
      <RouletteTableList />
    </div>
  );
}
