import { MyBets } from "@/components/roulette/my-bets";
import { PageHeader } from "@/components/shared/page-header";

export default function RouletteBetsPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="My Roulette Bets"
        description="Your betting history."
      />
      <MyBets />
    </div>
  );
}
