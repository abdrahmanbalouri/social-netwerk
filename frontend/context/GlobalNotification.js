"use client";
import { useEffect, useState } from "react";
import { useWS } from "./wsContext";
import Notification from "../components/notofication";
import { usePathname } from "next/navigation";
export default function GlobalNotification() {
  const { addListener, removeListener, connected } = useWS();
  const [toast, setToast] = useState(null);
  const pathname = usePathname();

  useEffect(() => {
    if (!connected) return;


    const handle = (data) => {
      
      const payload = data.data || data;
      switch (payload.subType) {
        case "follow":
          setToast(payload);
          break;
        case "unfollow":
          setToast(payload);
          break;
        case "followRequest":
          setToast(payload);
          break;
        case "message":
          if (pathname !== `/chat/${payload.from}`) {
            setToast(payload);
          }
          break;
        case "group_message":
          if (pathname !== `/groups/${payload.groupID}`) {
            setToast(payload);
          }
          break;
        case "group_invite":
          setToast(payload);
          break;
        case "group_join_request":
          setToast(payload);
          break;
        default:
          return;
      }




      setTimeout(() => setToast(null), 4500);
    };

    addListener("notification", handle);
    return () => removeListener("notification", handle);
  }, [connected, addListener, removeListener, pathname]);

  if (!toast) return null;

  return <Notification data={toast} />;
}