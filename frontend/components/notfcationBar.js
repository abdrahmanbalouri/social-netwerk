import { useState } from "react";
import "../styles/notfcationBar.css";
import formatTime from "../helpers/formatTime.js";
import { useRouter } from "next/navigation.js";


export default function NotBar({ notData }) {
    console.log("++++++++++++++");
    console.log(notData);
    
    console.log("++++++++++++++");
    
    const [filter, setFilter] = useState('all');
    const [notifications, setNotifications] = useState(notData || []);
    const router = useRouter();
    const Clear = () => {
        fetch('http://localhost:8080/api/clearNotifications', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",

        })

        setNotifications([]);
    };
    const getFilteredNotifications = () => {
        if (filter === 'unread') {
            return notifications.filter(noti => !noti.isRead);
        }
        return notifications;
    };

    const getNotificationIcon = (type) => {
        switch (type) {
            case 'like': return { icon: '❤️', className: 'icon-like' };
            case 'comment': return { icon: '💬', className: 'icon-comment' };
            case 'message': return { icon: '💬', className: 'icon-message' };
            default: return { icon: '🔔', className: 'icon-like' };
        }
    };
    const notificationtype = (noti) => {
        switch (noti.type) {
            case 'like':
            case 'comment':
            case 'message':
                router.push(`/chat/${noti.sender_id}`)
                break;
            case 'follow':
                router.push(`/profile/${noti.sender_id}`)
                break;
            case 'invite_to_group':
                router.push(`/profile/${noti.sender_id}`)
                break;
            case 'joinRequest':
                router.push(`/groups`)
                break;
            default:
                break;
        }
    }
    const filteredNotifications = getFilteredNotifications();

    return (
        <div className="dropdown-menu show">
            <div className="dropdown-header">
                <div className="dropdown-title">
                    Notifications
                </div>
                <div className="dropdown-tabs">
                </div>
            </div>

            <div className="notifications-list" id="notificationsList">
                {filteredNotifications.length === 0 ? (
                    <div className="empty-state">
                        <div className="empty-state-icon">🔔</div>
                        <p>No notifications available</p>
                    </div>
                ) : (
                    filteredNotifications.map((noti, index) => {
                        const iconData = getNotificationIcon(noti.type);
                        return (
                            <div
                                key={noti.id || index}
                                className={`notification-item ${!noti.isRead ? 'unread' : ''}`}
                                onClick={() => notificationtype(noti)}
                            >
                                <div className="notification-avatar" >
                                    <img
                                        src={noti?.photo ? `/uploads/${noti.photo}` : "/assets/default.png"}
                                        alt="user avatar"
                                        style={{
                                            width: '100%',
                                            height: '100%',
                                            objectFit: 'cover',
                                            borderRadius: '50%'
                                        }}
                                    />
                                    <div className={`notification-icon ${iconData.className}`}>
                                        {iconData.icon}
                                    </div>
                                </div>
                                <div className="notification-content">
                                    <div className="notification-text">
                                        <strong>{noti.name}</strong> {noti.message}
                                    </div>
                                    <div className="notification-time">{formatTime(noti.created_at)}</div>
                                </div>
                                {!noti.isRead && <div className="notification-dot"></div>}
                            </div>
                        );
                    })
                )}
            </div>

            <div className="dropdown-footer">
                <button className="see-all-button" onClick={Clear}>Clear</button>
            </div>
        </div>
    );
}