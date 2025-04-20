type Props = (
  | {
      srcLocation: string
      alt: string
      width: number
      height: number
      fill?: false
      priority?: false
    }
  | {
      srcLocation: string
      alt: string
      fill: true
      width?: never
      height?: never
      priority?: false
    }
) &
  React.ImgHTMLAttributes<HTMLImageElement>

export default function OptimizedImage({
  srcLocation,
  alt,
  width,
  height,
  fill,
  priority,
  ...props
}: Props) {
  return (
    <picture
      className={
        fill ? 'absolute inset-0 w-full h-full object-cover' : undefined
      }
    >
      <source
        media="(min-width: 1280px)"
        srcSet={`${srcLocation}_xl.webp`}
        type="image/webp"
      />

      <source
        media="(min-width: 1024px)"
        srcSet={`${srcLocation}_lg.webp`}
        type="image/webp"
      />

      <source
        media="(min-width: 768px)"
        srcSet={`${srcLocation}_md.webp`}
        type="image/webp"
      />

      <source
        media="(min-width: 640px)"
        srcSet={`${srcLocation}_sm.webp`}
        type="image/webp"
      />

      <source srcSet={`${srcLocation}.webp`} type="image/webp" />

      <img
        src={`${srcLocation}.png`}
        alt={alt}
        width={fill ? undefined : width}
        height={fill ? undefined : height}
        className={fill ? 'w-full h-full object-cover' : undefined}
        loading={priority ? 'eager' : 'lazy'}
        fetchPriority={priority ? 'high' : 'low'}
        {...props}
      />
    </picture>
  )
}
