import { PageHeader } from "@/components/shared/page-header";
import { BalanceCard } from "@/components/wallet/balance-card";
import { DepositForm } from "@/components/wallet/deposit-form";
import { TransactionList } from "@/components/wallet/transaction-list";
import { WithdrawForm } from "@/components/wallet/withdraw-form";

export default function WalletPage() {
  return (
    <div className="space-y-8">
      <PageHeader
        title="Wallet"
        description="Manage your funds and view transaction history."
      />
      <BalanceCard />
      <div className="grid gap-6 sm:grid-cols-2">
        <DepositForm />
        <WithdrawForm />
      </div>
      <TransactionList />
    </div>
  );
}
