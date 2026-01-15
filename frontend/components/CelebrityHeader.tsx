import Image from "next/image";
import { Celebrity, getEGOTStatus } from "@/lib/api";
import EGOTStatus from "./EGOTStatus";

interface Props {
  celebrity: Celebrity;
}

export default function CelebrityHeader({ celebrity }: Props) {
  const status = getEGOTStatus(celebrity.awards);

  return (
    <div className="flex flex-col md:flex-row items-center gap-8 mb-12">
      {/* Photo */}
      <div className="relative">
        <div className="photo-frame rounded-lg overflow-hidden bg-hollywood-charcoal">
          {celebrity.photo_url ? (
            <Image
              src={celebrity.photo_url}
              alt={celebrity.name}
              width={200}
              height={250}
              className="object-cover"
              unoptimized
            />
          ) : (
            <div className="w-[200px] h-[250px] flex items-center justify-center text-6xl text-gray-600">
              ★
            </div>
          )}
        </div>
      </div>

      {/* Info */}
      <div className="flex-1 text-center md:text-left">
        <h1 className="font-display text-4xl md:text-5xl font-bold text-hollywood-cream mb-4">
          {celebrity.name}
        </h1>
        <EGOTStatus status={status} size="md" />
        <p className="text-gray-500 text-sm mt-4">
          {status.count} of 4 awards achieved
        </p>
      </div>
    </div>
  );
}
