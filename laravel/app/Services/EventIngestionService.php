<?php

namespace App\Services;

use App\Models\Event;
use Carbon\Carbon;

class EventIngestionService
{
    /**
     * Ingest and persist an incoming application event synchronously.
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
}