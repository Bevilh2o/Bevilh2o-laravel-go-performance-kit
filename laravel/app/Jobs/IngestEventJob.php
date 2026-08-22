<?php

namespace App\Jobs;

use App\Services\EventIngestionService;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class IngestEventJob implements ShouldQueue
{
    use Queueable;

    /**
     * Create a new job instance.
     *
     * @param array{tenant: string, event: string, payload?: array|null, timestamp?: int|null} $eventData
     */
    public function __construct(
        public array $eventData
    ) {}

    /**
     * Execute the job.
     */
    public function handle(EventIngestionService $ingestionService): void
    {
        $ingestionService->ingestDirect($this->eventData);
    }
}