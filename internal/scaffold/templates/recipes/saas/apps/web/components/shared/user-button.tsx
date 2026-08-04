'use client';

import { User } from 'better-auth';

import { useSession } from '@repo/auth/client';
import { SidebarTrigger } from '@repo/ui/components/ui/sidebar';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@repo/ui/components/ui/tooltip';

export default function UserButton() {
  const { data, isPending } = useSession();

  if (isPending) return <div></div>;

  const user = data?.user;

  return (
    <div className="flex items-center justify-between gap-2 text-white">
      <div className="flex items-center gap-2 group-data-[collapsible=icon]:hidden">
        <div className="h-8 w-8 bg-blue-600 rounded-lg flex items-center justify-center font-bold">
          JD
        </div>
        <div>
          <h2 className="font-bold text-sm text-black">{user?.name}</h2>
          <p className="text-xs text-gray-500">{user?.email}</p>
        </div>
      </div>
      <Tooltip>
        <TooltipTrigger asChild>
          <SidebarTrigger className="text-gray-400 hover:text-white hover:bg-white/5 p-2" />
        </TooltipTrigger>
        <TooltipContent side={'right'}>
          <p>Abrir Menú</p>
        </TooltipContent>
      </Tooltip>
    </div>
  );
}
