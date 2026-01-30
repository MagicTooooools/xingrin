import type { ComponentProps, HTMLAttributes } from 'react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

export type StatusProps = ComponentProps<typeof Badge> & {
  status: 'online' | 'offline' | 'maintenance' | 'degraded';
};

// GitHub-style status color configuration
const statusStyles = {
  online: 'bg-[#238636]/10 text-[#238636] border-[#238636]/20 dark:text-[#3fb950]',
  offline: 'bg-[#da3633]/10 text-[#da3633] border-[#da3633]/20 dark:text-[#f85149]',
  maintenance: 'bg-[#848d97]/10 text-[#848d97] border-[#848d97]/20',
  degraded: 'bg-[#d29922]/10 text-[#d29922] border-[#d29922]/20',
};

export const Status = ({ className, status, ...props }: StatusProps) => (
  <Badge
    className={cn('flex items-center gap-2', 'group', status, statusStyles[status], className)}
    variant="outline"
    {...props}
  />
);

export type StatusIndicatorProps = HTMLAttributes<HTMLSpanElement>;

export const StatusIndicator = ({
  className,
  ...props
}: StatusIndicatorProps) => (
  <span className={cn('relative flex h-2 w-2', className)} {...props}>
    <span
      className={cn(
        'absolute inline-flex h-full w-full animate-ping rounded-full opacity-75',
        'group-[.online]:bg-[#238636] dark:group-[.online]:bg-[#3fb950]',
        'group-[.offline]:bg-[#da3633] dark:group-[.offline]:bg-[#f85149]',
        'group-[.maintenance]:bg-[#848d97]',
        'group-[.degraded]:bg-[#d29922]'
      )}
    />
    <span
      className={cn(
        'relative inline-flex h-2 w-2 rounded-full',
        'group-[.online]:bg-[#238636] dark:group-[.online]:bg-[#3fb950]',
        'group-[.offline]:bg-[#da3633] dark:group-[.offline]:bg-[#f85149]',
        'group-[.maintenance]:bg-[#848d97]',
        'group-[.degraded]:bg-[#d29922]'
      )}
    />
  </span>
);

export type StatusLabelProps = HTMLAttributes<HTMLSpanElement>;

export const StatusLabel = ({
  className,
  children,
  ...props
}: StatusLabelProps) => (
  <span className={cn(className)} {...props}>
    {children ?? (
      <>
        <span className="hidden group-[.online]:block">Online</span>
        <span className="hidden group-[.offline]:block">Offline</span>
        <span className="hidden group-[.maintenance]:block">Maintenance</span>
        <span className="hidden group-[.degraded]:block">Degraded</span>
      </>
    )}
  </span>
);
