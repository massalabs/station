import Intl from '@/i18n/i18n';

export function Loading() {
  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10 animate-pulse blur-sm">
          {Intl.t('search.title-banner')}
        </h1>
        <div className="overflow-auto h-[65vh]">
          <p className="mas-body mb-3 animate-pulse blur-sm">
            {Intl.t('search.fav-websites')}
          </p>
          <div className="flex flex-col gap-5 mb-14">
            <div>
              <div className="flex items-center space-x-4 animate-pulse">
                <div className="flex-1 py-1">
                  <div className="flex h-[60px] bg-tertiary rounded-lg mb-6 p-3.5 pl-3">
                    <div className="h-8 w-8 bg-f-disabled-2 rounded-full mr-2"></div>
                    <div>
                      <div className="h-2 w-40 bg-f-disabled-2 rounded mt-1.5"></div>
                      <div className="h-1 w-20 bg-f-disabled-2 rounded mt-1.5"></div>
                    </div>
                  </div>

                  <div className="flex h-[60px] bg-tertiary rounded-lg mb-6 p-3.5 pl-3">
                    <div className="h-8 w-8 bg-f-disabled-2 rounded-full mr-2"></div>
                    <div>
                      <div className="h-2 w-36 bg-f-disabled-2 rounded mt-1.5"></div>
                      <div className="h-1 w-20 bg-f-disabled-2 rounded mt-1.5"></div>
                    </div>
                  </div>

                  <div className="flex h-[60px] bg-tertiary rounded-lg mb-6 p-3.5 pl-3">
                    <div className="h-8 w-8 bg-f-disabled-2 rounded-full mr-2"></div>
                    <div>
                      <div className="h-2 w-48 bg-f-disabled-2 rounded mt-1.5"></div>
                      <div className="h-1 w-20 bg-f-disabled-2 rounded mt-1.5"></div>
                    </div>
                  </div>

                  <div className="flex h-[60px] bg-tertiary rounded-lg mb-6 p-3.5 pl-3">
                    <div className="h-8 w-8 bg-f-disabled-2 rounded-full mr-2"></div>
                    <div>
                      <div className="h-2 w-36 bg-f-disabled-2 rounded mt-1.5"></div>
                      <div className="h-1 w-20 bg-f-disabled-2 rounded mt-1.5"></div>
                    </div>
                  </div>

                  <div className="flex h-[60px] bg-tertiary rounded-lg mb-6 p-3.5 pl-3">
                    <div className="h-8 w-8 bg-f-disabled-2 rounded-full mr-2"></div>
                    <div>
                      <div className="h-2 w-40 bg-f-disabled-2 rounded mt-1.5"></div>
                      <div className="h-1 w-20 bg-f-disabled-2 rounded mt-1.5"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
