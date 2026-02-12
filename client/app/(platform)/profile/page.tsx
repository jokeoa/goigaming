import { ProfileCard } from "@/components/profile/profile-card";
import { ProfileStats } from "@/components/profile/profile-stats";
import { PageHeader } from "@/components/shared/page-header";

export default function ProfilePage() {
  return (
    <div className="space-y-8">
      <PageHeader title="Profile" description="Your account details." />
      <div className="grid gap-6 lg:grid-cols-2">
        <ProfileCard />
        <ProfileStats />
      </div>
    </div>
  );
}
