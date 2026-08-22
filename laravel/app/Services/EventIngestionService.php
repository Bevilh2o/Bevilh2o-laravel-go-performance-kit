<?php

namespace App\Services;

use App\Jobs\IngestEventJob;
use App\Models\Event;
use Carbon\Carbon;

class EventIngestionService
{
    /**
     * Ingest and persist an incoming application event synchronously into PostgreSQL.
     *
     * @param array{tenant: string, event: string, payload?: array|null, timestamp?: int|null} $data
     * @return Event
     */
    public function ingestDirect(array $data): Event
    {
        $occurredAt = isset($data['timestamp'])
            ? Carbon::createFromTimestamp($data['timestamp'])
            : now();

        return Event::create([
            'tenant_id' => $data['tenant'],
            'event_type' => $data['event'],
            'payload' => $data['payload'] ?? null,
            'occurred_at' => $occurredAt,
        ]);
    }

    /**
     * Push event to Redis queue for asynchronous background processing.
     *
     * @param array{tenant: string, event: string, payload?: array|null, timestamp?: int|null} $data
     * @return void
     */
    public function ingestAsync(array $data): void
    {
        IngestEventJob::dispatch($data);
    }
}