import { useState } from "react";
import "../styles/notfcationBar.css";

export default function NotBar({ notData }) {
    const [filter, setFilter] = useState('all');
    const [notifications, setNotifications] = useState(notData || []);


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

    const filteredNotifications = getFilteredNotifications();

    return (
        <div className="notification_dd">
            <div className="dropdown-menu show">
                <div className="dropdown-header">
                    <div className="dropdown-title">
                        Notifications
                        <button className="settings-button">⚙️</button>
                    </div>
                    <div className="dropdown-tabs">
                        <button
                            className={`tab-button ${filter === 'all' ? 'active' : ''}`}
                        >
                            All
                        </button>
                       
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
                                    onClick={() => markAsRead(index)}
                                >
                                    <div className="notification-avatar">
                                        <img
                                            src={noti?.photo ? `/uploads/${noti.photo}` : "/uploads/default.png"}
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
                                        <div className="notification-time">{noti.created_at}</div>
                                    </div>
                                    {!noti.isRead && <div className="notification-dot"></div>}
                                </div>
                            );
                        })
                    )}
                </div>

                <div className="dropdown-footer">
                    <button className="see-all-button">Show All Notifications</button>
                </div>
            </div>
        </div>
    );
}