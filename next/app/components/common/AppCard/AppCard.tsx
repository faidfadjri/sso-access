import { AppCardProps } from "./AppCard.type";
import Image from "next/image";

function AppCard({ name, iconPath, className = "", onClick }: AppCardProps) {
  return (
    <div onClick={onClick} className={`flex flex-col items-center justify-center bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow cursor-pointer py-8 px-4 w-full h-[140px] border border-gray-100 ${className}`}>
      <div className={`mb-3 p-3 rounded-full bg-blue-50/50`}>
        <Image src={iconPath || "/apps/website.png"} alt={name} width={40} height={40} unoptimized/>
      </div>
      <span className="font-semibold text-gray-700 text-sm md:text-base text-center">{name}</span>
    </div>
  );
}

function AppCardSkeleton({ className = "" }: { className?: string }) {
  return (
    <div className={`flex flex-col items-center justify-center bg-white rounded-xl shadow-sm py-8 px-4 w-full h-[140px] border border-gray-100 animate-pulse ${className}`}>
      <div className="mb-3 p-3 rounded-full bg-gray-100 w-[64px] h-[64px]"></div>
      <div className="h-4 bg-gray-100 rounded w-20"></div>
    </div>
  );
}

AppCard.Skeleton = AppCardSkeleton;

export default AppCard;
