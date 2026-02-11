import { QuickActions } from "@/components/dashboard/quick-actions";
import { StatsCards } from "@/components/dashboard/stats-cards";
import { TableList } from "@/components/dashboard/table-list";
import { PageHeader } from "@/components/shared/page-header";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <PageHeader
        title="Dashboard"
        description="Welcome back. Here's your overview."
      />
      <StatsCards />
      <QuickActions />
      <div className="space-y-4">
        <h2 className="text-lg font-semibold tracking-tight">Open Tables</h2>
        <TableList />
      </div>
    </div>
  );
}
