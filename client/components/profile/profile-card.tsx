"use client";

import { AvatarDisplay } from "@/components/profile/avatar-display";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { useCurrentUser } from "@/hooks/use-user";
import { formatDate } from "@/lib/utils";

export function ProfileCard() {
  const { data: user, isLoading } = useCurrentUser();

  if (isLoading) {
    return (
      <Card className="border-border">
        <CardContent className="flex items-center gap-6 p-6">
          <Skeleton className="h-20 w-20 rounded-lg" />
          <div className="space-y-2">
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-4 w-48" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!user) return null;

  return (
    <Card className="border-border">
      <CardHeader>
        <CardTitle className="text-sm">Profile</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-6">
          <AvatarDisplay seed={user.id} size={80} />
          <div>
            <h2 className="text-lg font-semibold">{user.username}</h2>
            <p className="text-sm text-muted-foreground">{user.email}</p>
          </div>
        </div>
        <Separator />
        <div className="grid gap-3 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">User ID</span>
            <span className="font-mono text-xs">{user.id}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Member since</span>
            <span>{formatDate(user.created_at)}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
