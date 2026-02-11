import { ArrowDownToLine, ArrowUpFromLine, User } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function QuickActions() {
  return (
    <Card className="border-border">
      <CardHeader className="pb-3">
        <CardTitle className="text-sm">Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-wrap gap-2">
        <Button asChild variant="outline" size="sm">
          <Link href="/wallet">
            <ArrowDownToLine className="mr-2 h-3.5 w-3.5" />
            Deposit
          </Link>
        </Button>
        <Button asChild variant="outline" size="sm">
          <Link href="/wallet">
            <ArrowUpFromLine className="mr-2 h-3.5 w-3.5" />
            Withdraw
          </Link>
        </Button>
        <Button asChild variant="outline" size="sm">
          <Link href="/profile">
            <User className="mr-2 h-3.5 w-3.5" />
            Profile
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
