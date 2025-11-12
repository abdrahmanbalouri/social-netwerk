'use client';
import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { useWS } from '../context/wsContext'; // adjust the path

export default function WSRouteHandler() {
  const pathname = usePathname();
  const { disconnect } = useWS();

  useEffect(() => {
    // Close WS if user navigates to login/register
    if (pathname === '/login' || pathname === '/register') {
      disconnect();
    }
  }, [pathname, disconnect]);

  return null;
}
