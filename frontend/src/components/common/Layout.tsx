import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen animate-gradient bg-gradient-to-br from-rose-100 via-violet-100 to-teal-100">
      <main className="flex-1 flex items-center justify-center min-h-screen w-full">
        {children}
      </main>
    </div>
  );
}

export default Layout; 