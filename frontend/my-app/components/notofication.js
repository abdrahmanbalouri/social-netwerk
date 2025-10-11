import Link from "next/link";
import "../styles/notifaction.css";

export default function Notification({ data }) {
    if (!data) return null;
    console.log();


    return (
        <div className="chat">
            <Link href={`/profile/${data.receiverId}`}>
                <div className="Profile">
                    <img src={data.image || "/uploads/default.png"} alt={data.name} />
                </div>
                <div className="message">
                    You have a new follower
                </div>
                <div className="User">
                    {data.name}
                </div>
            </Link>
        </div>
    );
}
