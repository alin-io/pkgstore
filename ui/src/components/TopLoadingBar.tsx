import { Container } from './Container';

interface Props {
  showSkeleton?: boolean;
}

export function TopLoadingBar({ showSkeleton }: Props) {
  return (
    <>
      <div className="fixed top-16 left-0 w-full h-1 bg-gray-200 z-50">
        <div className="h-full bg-blue-400 animate-pulse"></div>
      </div>
      {showSkeleton && (
        <Container className="flex justify-center">
          <div role="status" className="space-y-2.5 animate-pulse w-full">
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-200 rounded-full w-32"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-24"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
            </div>
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-200 rounded-full w-24"></div>
            </div>
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-200 rounded-full w-80"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
            </div>
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-200 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-24"></div>
            </div>
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-300 rounded-full w-32"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-24"></div>
              <div className="h-2.5 bg-gray-200 rounded-full w-full"></div>
            </div>
            <div className="flex items-center w-full space-x-2">
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
              <div className="h-2.5 bg-gray-200 rounded-full w-80"></div>
              <div className="h-2.5 bg-gray-300 rounded-full w-full"></div>
            </div>
            <span className="sr-only">Loading...</span>
          </div>
        </Container>
      )}
    </>
  );
}
