import { SidebarProvider } from '@repo/ui/components/ui/sidebar';

import { AppSidebar } from '../../../components/dashboard/app-sidebar';

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        <AppSidebar />

        {/* Main Content Area - "Encerrado" */}
        <main className="flex-1 p-4 h-full">
          <div className="h-full w-full rounded-3xl bg-white overflow-hidden border border-gray-200 relative flex flex-col ">
            {children}
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
}
