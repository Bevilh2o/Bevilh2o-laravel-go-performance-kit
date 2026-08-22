<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\StoreEventRequest;
use App\Services\EventIngestionService;
use Illuminate\Http\JsonResponse;
use Symfony\Component\HttpFoundation\Response;

class EventController extends Controller
{
    public function __construct(
        protected EventIngestionService $ingestionService
    ) {}

    /**
     * Synchronous baseline endpoint (HTTP -> DB).
     */
    public function store(StoreEventRequest $request): JsonResponse
    {
        $event = $this->ingestionService->ingestDirect($request->validated());

        return response()->json([
            'status' => 'success',
            'data' => [
                'id' => $event->id,
                'tenant_id' => $event->tenant_id,
                'event_type' => $event->event_type,
                'occurred_at' => $event->occurred_at->toIso8601String(),
            ],
        ], Response::HTTP_CREATED);
    }

    /**
     * Asynchronous baseline endpoint (HTTP -> Redis Queue -> 202 Accepted).
     */
    public function storeAsync(StoreEventRequest $request): JsonResponse
    {
        $this->ingestionService->ingestAsync($request->validated());

        return response()->json([
            'status' => 'queued',
            'message' => 'Event accepted for background processing',
        ], Response::HTTP_ACCEPTED);
    }
}