import OptimizedImage from "@/utils/image";

export default function Features() {
  return (
    <div className="flex flex-col h-fit  items-center justify-start gap-10 w-full">
      <div className="flex flex-col items-center justify-start gap-10 w-full">
        <h2 className="font-semibold text-4xl text-center md:text-5xl tracking-large text-brand-text">
          Features
        </h2>
      </div>

      <div className="h-fit md:grid-cols-2 grid-cols-1 w-full gap-10 grid">
        <FeatureCard
          title="Extensions"
          description="Save vocabulary with ease on both Firefox and Chrome extensions."
          imageSrc="extensions"
        />

        <FeatureCard
          title="Offline Support"
          description="Spotty connection? Vocab Thing works completely offline with all your data"
          imageSrc="offline"
        />

        <FeatureCard
          title="AI Selection"
          description="Not sure the best phrase for a convo? Let AI handle it for you"
          imageSrc="ai"
        />

        <FeatureCard
          title="Profile"
          description="Create your own profile, with vocabulary you choose to showcase"
          imageSrc="profile"
        />
      </div>
    </div>
  );
}

function FeatureCard({
  title,
  description,
  imageSrc,
}: {
  title: string;
  description: string;
  imageSrc: string;
}) {
  return (
    <div className="h-[300px] flex-col  relative bg-brand-primary-dark rounded-lg w-full flex  gap-5 items-start justify-end py-9 px-10">
      <div className=" absolute top-9 right-10 size-20 md:size-24 xl:size-28">
        <OptimizedImage
          srcLocation={`/home/features/${imageSrc}/${imageSrc}`}
          alt={title}
          sizes="(min-width: 1280px) 112px, (min-width: 768px) 96px, 80px"
          fill
          className="object-cover"
        />
      </div>
      <h3 className="text-brand-text text-3xl tracking-large font-medium">
        {title}
      </h3>
      <p className="text-brand-text text-base w-[80%] tracking-small">
        {description}
      </p>
    </div>
  );
}
