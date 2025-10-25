"use client";
import { useEffect, useState } from "react";
import { useWS } from "./wsContext";
import Notification from "../components/notofication";
import { useDarkMode } from "./darkMod";
import { usePathname } from "next/navigation"; 
export default function GlobalNotification() {
  const { addListener, removeListener, connected } = useWS();
  const { darkMode } = useDarkMode();
  const [toast, setToast] = useState(null);
  const pathname = usePathname(); 

  useEffect(() => {
    if (!connected) return;

    const handle = (data) => {
      const payload = data.data || data;
console.log("Notification received in GlobalNotification:", payload);

      if (payload.type!=="message"||pathname !== `/chat/${payload.from}`) {
        setToast(payload);
      }

      setTimeout(() => setToast(null), 4500); 
    };

    addListener("notification", handle);
    return () => removeListener("notification", handle);
  }, [connected, addListener, removeListener, pathname]);

  if (!toast) return null;

  return <Notification data={toast} />;
}
