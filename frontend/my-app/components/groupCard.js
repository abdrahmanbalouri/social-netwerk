import { Users, ChevronRight } from "lucide-react";

export function GroupCard({ group, onShow }) {
  console.log("inside groupCard function/component");
  return (
    <div className="group-card">
      <div className="group-header">
        <div className="group-icon-wrapper">
          <div className="group-icon">
            <Users />
          </div>
        </div>
      </div>
      {/* Group Content */}
      <div className="group-content">
        <h2 className="group-title">{group.Title}</h2>
        <p className="group-description">{group.Description}</p>
        {/* Group Footer */}
        <div className="group-footer">
          <div className="group-members">
            <Users />
            <span className="members-text-full">{group.MemberCount} members</span>
            <span className="members-text-short">{group.MemberCount}</span>
          </div>
          <button className="view-button" onClick={() => onShow(group.ID)}>
            <span>View</span>
            <ChevronRight />
          </button>
        </div>
      </div>
    </div>
  );
}
