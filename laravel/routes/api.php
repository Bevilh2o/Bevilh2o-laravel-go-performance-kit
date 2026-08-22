<?php

use App\Http\Controllers\Api\EventController;
use Illuminate\Support\Facades\Route;

// Synchronous baseline (direct DB insertion)
Route::post('/events', [EventController::class, 'store']);

// Asynchronous baseline (queued to Redis)
Route::post('/events/async', [EventController::class, 'storeAsync']);