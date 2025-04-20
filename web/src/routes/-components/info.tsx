import { Button } from '@/components/ui/button'
import OptimizedImage from '@/lib/image'

export default function Info() {
  return (
    <div className="w-full flex flex-col md:flex-row items-center justify-between gap-10">
      <div className="w-full flex flex-col h-fit  items-center md:items-start justify-start gap-10 md:w-[80%] xl:w-[70%]">
        <h2 className="font-semibold text-4xl text-center md:text-5xl tracking-large text-brand-text">
          Info
        </h2>

        <p className="text-brand-text font-medium text-2xl text-center md:text-left">
          Hiya! I am Arinjay Dhar. A self proclaimed “indie hacker” trying to
          make stuff not explode.
        </p>

        <p className="text-center md:text-left text-brand-text font-medium text-2xl">
          Like what you see here? Feel free to check out my{' '}
          <span className="text-brand-accent">
            <a
              href="https://arinji.com"
              target="_blank"
              rel="noopener noreferrer"
            >
              portfolio
            </a>
          </span>{' '}
          for all my other projects.
        </p>

        <p className="text-brand-text font-medium text-2xl text-center md:text-left">
          I run Vocab Thing on my own costs, but if you like it and want to
          support me, feel free to buy me a coffee :D
        </p>
        <div className="w-full h-fit flex flex-row items-center justify-center md:justify-start gap-10">
          <Button variant={'default'} size={'lg'} asChild>
            <a
              href="https://github.com/Arinji2/vocab-thing"
              target="_blank"
              rel="noopener noreferrer"
            >
              Vocab Thing Monorepo
            </a>
          </Button>
          <Button variant={'secondary'} size={'lg'} asChild>
            <a
              href="https://buymeacoffee.com/arinjii"
              target="_blank"
              rel="noopener noreferrer"
            >
              Buy Me A Coffee
            </a>
          </Button>
        </div>
      </div>
      <div className="relative shrink-0 rounded-full overflow-hidden xl:size-44 md:size-40 size-36">
        <OptimizedImage
          srcLocation="/home/info/info"
          alt="Arinji Profile"
          fill
        />
      </div>
    </div>
  )
}
